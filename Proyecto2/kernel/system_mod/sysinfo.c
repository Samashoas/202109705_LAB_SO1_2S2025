#include <linux/init.h>
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/proc_fs.h>
#include <linux/seq_file.h>
#include <linux/mm.h>
#include <linux/sched/signal.h>
#include <linux/sched.h>
#include <linux/sysinfo.h>

# define PROC_NAME "sysinfo_so1_202109705"

MODULE_LICENSE("GPL");
MODULE_AUTHOR("Juan Pablo Samayoa Ruiz-202109705");
MODULE_DESCRIPTION("SO1 - MONITOREO DE PROCESOS Y SISTEMA");
MODULE_VERSION("0.1");

static char state_char(const struct task_struct *t){
    if (task_is_running(t)) return 'R';
    if (READ_ONCE(t->__state) == TASK_INTERRUPTIBLE) return 'S';
    if (READ_ONCE(t->__state) == TASK_UNINTERRUPTIBLE) return 'D';
    if (READ_ONCE(t->__state) & __TASK_STOPPED) return 'T';
    if (READ_ONCE(t->__state) & __TASK_TRACED) return 't';
    if (READ_ONCE(t->__state) & EXIT_ZOMBIE) return 'Z';
    if (READ_ONCE(t->__state) & EXIT_DEAD) return 'X';
    return '?';
}

static int sysinfo_show(struct seq_file *m, void *v){
    struct sysinfo si;
    si_meminfo(&si);

    unsigned long total_kb = (si.totalram * (unsigned long)PAGE_SIZE) >> 10;
    unsigned long free_kb = (si.freeram * (unsigned long)PAGE_SIZE) >> 10;
    unsigned long used_kb = (total_kb > free_kb) ? (total_kb - free_kb) : 0;

    seq_puts(m, "{\n");
    seq_printf(m, "\"ram\" : {\"total_kb\": %lu, \"used_kb\": %lu, \"free_kb\": %lu},\n", total_kb, used_kb, free_kb);
    seq_puts(m, "\"processes\" : [\n");

    {
        struct task_struct *task;
        bool first = true;

        for_each_process(task){
            if(!first) seq_puts(m, ",\n");
            first = false;
            seq_printf(m, "{\"pid\": %d, \"state\": \"%c\"}", task->pid, state_char(task));
        }
    }

    seq_puts(m, "\n]\n}\n");
    return 0;
}

static int sysinfo_open(struct inode *inode, struct file *file){
    return single_open(file, sysinfo_show, NULL);
}

static const struct proc_ops fops ={
    .proc_open = sysinfo_open,
    .proc_read = seq_read,
    .proc_lseek = seq_lseek,
    .proc_release = single_release,
};

static int __init sysinfo_init(void){
    if(!proc_create(PROC_NAME, 0444, NULL, &fops)){
        pr_err("No se pudo crear el proc %s\n", PROC_NAME);
        return -ENOMEM;
    }

    pr_info("sysinfo_so1 cargado\n");
    return 0;
}

static void __exit sysinfo_exit(void){
    remove_proc_entry(PROC_NAME, NULL);
    pr_info("sysinfo_so1 descargado\n");
}

module_init(sysinfo_init);
module_exit(sysinfo_exit);