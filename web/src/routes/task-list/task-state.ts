import { Task } from "./task.model";

export interface TaskState {
    initialLoad: boolean;
    tasks: Task[];
    newTaskName: string;
}

