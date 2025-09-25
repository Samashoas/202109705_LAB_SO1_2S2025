#include <linux/init.h>
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/proc_fs.h>
#include <linux/seq_file.h>
#include <linux/mm.h>
#include <linux/sched/signal.h>
#include <linux/sched.h>
#include <linux/sysinfo.h>
#include <linux/cgroup.h>
#include <linux/nsproxy.h>
#include <linux/fs.h>
#include <linux/uaccess.h>
#include <linux/pid_namespace.h>

#define PROC_NAME "continfo_so1_202109705"

MODULE_LICENSE("GPL");
MODULE_AUTHOR("Juan Pablo Samayoa Ruiz-202109705");
MODULE_DESCRIPTION("Modulo del kernel que muestra informacion de contenedores Docker");
MODULE_VERSION("0.1");

/* Función para verificar si un proceso pertenece a un contenedor Docker */
static bool is_docker_container(struct task_struct *task) {
    /* Si no tiene un espacio de nombres PID separado, no es un contenedor */
    if (!task->nsproxy || !task->nsproxy->pid_ns_for_children)
        return false;

    /* Si está en el namespace PID inicial, no es un contenedor */
    if (task->nsproxy->pid_ns_for_children == &init_pid_ns)
        return false;

    /* Si tiene un nivel de namespace diferente, probablemente es un contenedor */
    if (task->nsproxy->pid_ns_for_children->level != 0)
        return true;

    return false;
}

/* Obtiene información de memoria del proceso */
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

/* Calcula porcentajes aproximados de CPU y memoria */
static void calc_cpu_mem_usage(const struct task_struct *task, 
                              unsigned long total_ram, 
                              unsigned long *cpu_usage, 
                              unsigned long *mem_usage)
{
    unsigned long utime, stime;
    unsigned long long runtime;
    
    /* Cálculo simplificado de uso de CPU */
    utime = task->utime;
    stime = task->stime;
    
    runtime = utime + stime;
    *cpu_usage = runtime ? (runtime * 10) / (jiffies + 1) : 0;
    if (*cpu_usage > 100) *cpu_usage = 100;
    
    /* Cálculo de uso de memoria */
    if (total_ram > 0 && task->mm) {
        *mem_usage = task->mm->total_vm * 100 / (total_ram << (PAGE_SHIFT - 10));
        if (*mem_usage > 100) *mem_usage = 100;
    } else {
        *mem_usage = 0;
    }
}

static int continfo_show(struct seq_file *m, void *v){
    struct sysinfo si;
    unsigned long total_kb;
    
    /* Obtener información de memoria total */
    si_meminfo(&si);
    total_kb = (si.totalram * (unsigned long)PAGE_SIZE) >> 10;
    
    seq_puts(m, "{\n \"container_processes\": [\n");

    {
        struct task_struct *task;
        bool first = true;
        unsigned long vsz_kb, rss_kb, cpu_pct, mem_pct;

        /* Recorrer todos los procesos */
        for_each_process(task) {
            /* Solo incluir procesos que parecen ser de contenedores */
            if (is_docker_container(task)) {
                get_memory_kb(task, &vsz_kb, &rss_kb);
                calc_cpu_mem_usage(task, total_kb, &cpu_pct, &mem_pct);
                
                if(!first) seq_puts(m, ",\n");
                first = false;

                seq_puts(m, "  {\n");
                seq_printf(m, "    \"pid\": %d,\n", task->pid);
                seq_printf(m, "    \"name\": \"%s\",\n", task->comm);
                seq_printf(m, "    \"container_id\": \"%d\",\n", task->nsproxy->pid_ns_for_children->level);
                seq_printf(m, "    \"cmd\": \"%s\",\n", task->comm);
                seq_printf(m, "    \"vsz_kb\": %lu,\n", vsz_kb);
                seq_printf(m, "    \"rss_kb\": %lu,\n", rss_kb);
                seq_printf(m, "    \"pct_mem\": %lu,\n", mem_pct);
                seq_printf(m, "    \"pct_cpu\": %lu\n", cpu_pct);
                seq_puts(m, "  }");
            }
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