#include <linux/init.h>
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/proc_fs.h>
#include <linux/seq_file.h>
#include <linux/mm.h>
#include <linux/sched/signal.h>
#include <linux/sched.h>
#include <linux/sysinfo.h>
#include <linux/nsproxy.h>
#include <linux/pid_namespace.h>

#define PROC_NAME "sysinfo_so1_202109705"

MODULE_LICENSE("GPL");
MODULE_AUTHOR("Juan Pablo Samayoa Ruiz-202109705");
MODULE_DESCRIPTION("SO1 - MONITOREO DE PROCESOS Y SISTEMA");
MODULE_VERSION("0.1");

/* Función para obtener el estado del proceso como un caracter - compatible con kernel 6.14 */
static char state_char(const struct task_struct *tsk)
{
    /* En kernels 6.x+ necesitamos usar task_state_to_char() o impl. propia */
    unsigned long state;
    
    state = READ_ONCE(tsk->__state);
    
    if (tsk->exit_state & EXIT_DEAD)
        return 'X';
    if (tsk->exit_state & EXIT_ZOMBIE)
        return 'Z';
    if (state & TASK_UNINTERRUPTIBLE)
        return 'D';
    if (state & TASK_STOPPED)
        return 'T';
    if (state & TASK_TRACED)
        return 't';
    if (state == TASK_INTERRUPTIBLE)
        return 'S';
    
    return 'R';
}

/* Verifica si un proceso es de contenedor */
static bool is_container_process(const struct task_struct *task)
{
    /* Un proceso de contenedor típicamente tiene un namespace PID distinto */
    if (task->nsproxy && 
        task->nsproxy->pid_ns_for_children && 
        task->nsproxy->pid_ns_for_children != &init_pid_ns) {
        return true;
    }
    
    return false;
}

/* Implementación para mostrar la información en /proc */
static int sysinfo_show(struct seq_file *m, void *v)
{
    struct sysinfo si;
    
    /* Obtener información de memoria */
    si_meminfo(&si);
    
    unsigned long total_kb = (si.totalram * (unsigned long)PAGE_SIZE) >> 10;
    unsigned long free_kb = (si.freeram * (unsigned long)PAGE_SIZE) >> 10;
    unsigned long used_kb = total_kb - free_kb;
    
    /* Escribir información de RAM en formato JSON */
    seq_puts(m, "{\n");
    seq_printf(m, "  \"ram\": {\n");
    seq_printf(m, "    \"total_kb\": %lu,\n", total_kb);
    seq_printf(m, "    \"used_kb\": %lu,\n", used_kb);
    seq_printf(m, "    \"free_kb\": %lu\n", free_kb);
    seq_printf(m, "  }\n"); /* Nota: eliminé la coma al final */
    seq_puts(m, "}\n");
    
    return 0;
}

/* Función de apertura del archivo /proc */
static int sysinfo_open(struct inode *inode, struct file *file)
{
    return single_open(file, sysinfo_show, NULL);
}

/* Operaciones del archivo /proc (formato moderno) */
static const struct proc_ops fops = {
    .proc_open = sysinfo_open,
    .proc_read = seq_read,
    .proc_lseek = seq_lseek,
    .proc_release = single_release,
};

/* Función de inicialización del módulo */
static int __init sysinfo_init(void)
{
    /* Crear la entrada en /proc */
    if (!proc_create(PROC_NAME, 0444, NULL, &fops)) {
        pr_err("No se pudo crear el archivo /proc/%s\n", PROC_NAME);
        return -ENOMEM;
    }
    
    pr_info("Módulo sysinfo_so1_202109705 cargado correctamente\n");
    return 0;
}

/* Función de limpieza del módulo */
static void __exit sysinfo_exit(void)
{
    /* Eliminar la entrada de /proc */
    remove_proc_entry(PROC_NAME, NULL);
    pr_info("Módulo sysinfo_so1_202109705 descargado correctamente\n");
}

/* Registrar las funciones de inicialización y salida */
module_init(sysinfo_init);
module_exit(sysinfo_exit);