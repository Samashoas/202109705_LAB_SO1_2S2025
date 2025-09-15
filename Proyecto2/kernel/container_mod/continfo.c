#include <linux/init.h>
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/proc_fs.h>
#include <linux/seq_file.h>
#include <linux/mm.h>
#include <linux/sched/signal.h>
#include <linux/sched.h>
#include <linux/sysinfo.h>

#define PROC_NAME "continfo_so1_202109705"

MODULE_LICENSE("GPL");
MODULE_AUTHOR("Juan Pablo Samayoa Ruiz-202109705");
MODULE_DESCRIPTION("Modulo del kernel que muestra informacion de contenedores");
MODULE_VERSION("0.1");

static void get_memory_kb(struct task_struct *t, unsigned long *vsz_kb, unsigned long *rss_kb){
    *vsz_kb = 0;
    *rss_kb = 0;

    if(t->mm){
        unsigned long pages_vsz = t->mm->total_vm;
        unsigned long pages_rss = get_mm_rss(t->mm);
        *vsz_kb = (pages_vsz * (unsigned long)PAGE_SIZE) >> 10;
        *rss_kb = (pages_rss * (unsigned long)PAGE_SIZE) >> 10;
    }
}

static int continfo_show(struct seq_file *m, void *v){
    seq_puts(m, "{\n \"containers\": [\n");

    {
        struct task_struct *task;
        bool first = true;
        unsigned long vsz_kb, rss_kb;

        for_each_process(task){
            get_memory_kb(task, &vsz_kb, &rss_kb);
            if(!first) seq_puts(m, ",\n");
            first = false;

            seq_printf(m, "{\"pid\": %d, \"name\": \"%s\", \"cmd\": \"\", \"vsz_kb\": %lu, \"rss_kb\": %lu, \"pct_mem\": 0, \"pct_cpu\": 0}", 
                task->pid, task->comm, vsz_kb, rss_kb);
        }
    }
    seq_puts(m, "\n ]\n}\n");
    return 0;
}

static int continfo_open(struct inode *inode, struct file *file){
    return single_open(file, continfo_show, NULL);
}

static const struct proc_ops fops = {
    .proc_open = continfo_open,
    .proc_read = seq_read,
    .proc_lseek = seq_lseek,
    .proc_release = single_release,
};

static int __init continfo_init(void){
    if(!proc_create(PROC_NAME, 0444, NULL, &fops)){
        pr_err("No se pudo crear el archivo /proc/%s\n", PROC_NAME);
        return -ENOMEM;
    }
    pr_info("/proc/%s creado\n", PROC_NAME);
    return 0;
}

static void __exit continfo_exit(void){
    remove_proc_entry(PROC_NAME, NULL);
    pr_info("/proc/%s removido\n", PROC_NAME);
}

module_init(continfo_init);
module_exit(continfo_exit);