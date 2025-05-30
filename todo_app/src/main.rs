use std::env;

struct Todo {
    tasks: Vec<String>,
}

impl Todo {
    fn add(&mut self, task: String) {
        self.tasks.push(task.clone());
        println!("Added new task: {task}");
    }
    fn done(&mut self, index: usize) {
        self.tasks.remove(index);
        println!("Removed a task");
    }
    fn list(&self) {
        for (i, task) in self.tasks.iter().enumerate() {
            println!("{}, {}", i + 1, task);
        }
    }
}
fn main() {
    let mut todo = Todo { tasks: vec![] };

    let args: Vec<String> = env::args().collect();
    let command = args.get(1).map(|s| s.as_str());

    match command {
        Some("add") => {
            let task = args[2].clone();
            todo.add(task);
        }
        Some("done") => {
            if let Ok(idx) = args[2].parse::<usize>() {
                todo.done(idx);
            };
        }
        Some("list") => {
            todo.list();
        }
        _ => {
            todo!()
        }
    }
}
