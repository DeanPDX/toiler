import { Task } from "../routes/task-list/task.model";
import { SessionService } from "./session.service";

/**
 * A simple service to wrap our API calls and return data.
 */
export class APIService {
    readonly baseURL = '/api';

    /**
     * HTTP get using fetch API.
     * @param url The URL of the API route you want to get.
     */
    private get<T>(url: string): Promise<T> {
        return fetch(url,
            {
                method: 'get',
                headers: {
                    'Accept': 'application/json, text/plain, */*',
                    'Content-Type': 'application/json',
                    'X-Auth': SessionService.getAuthToken()
                }
            })
            .then(response => {
                if (!response.ok) {
                    throw new Error(`HTTP error: ${response.status}`)
                }
                return response.json()
            });
    }

    /**
     * HTTP post using fetch API.
     * @param url The URL you're posting to
     * @param body The body of the request
     */
    private post<ResponseType>(url: string, body: any): Promise<ResponseType> {
        return fetch(url, {
            method: 'post',
            headers: {
                'Accept': 'application/json, text/plain, */*',
                'Content-Type': 'application/json',
                'X-Auth': SessionService.getAuthToken()
            },
            body: JSON.stringify(body)
        }).then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error: ${response.status}`)
            }
            return response.json()
        });
    }

    createAccount(email: string, password: string): Promise<boolean> {
        return this.post<authResult>(`${this.baseURL}/createAccount`, { email: email, password: password }).then(value => {
            SessionService.setAuthToken(value.token);
            return value.success;
        })
    }

    logIn(email: string, password: string): Promise<boolean> {
        return this.post<authResult>(`${this.baseURL}/authenticate`, { email: email, password: password }).then(value => {
            SessionService.setAuthToken(value.token);
            return value.success;
        })
    }

    /**
     * Get the current task list items.
     */
    getTaskList(): Promise<Task[]> {
        return this.get<Task[]>(`${this.baseURL}/tasks/list`)
    }

    getTaskListExcel() {
        return fetch('/api/tasks/excel', {
          headers: {
            'X-Auth': SessionService.getAuthToken(),
          },
        })
          .then(res => res.blob())
      }

    /**
     * Create a new task and return true for success.
     * @param taskName The name of the new task
     */
    addTask(taskName: string): Promise<boolean> {
        return this.post<boolean>(`${this.baseURL}/tasks/add`, { taskName: taskName });
    }

    /**
     * Update a task.
     * @param taskID The task ID you want to update
     * @param completed Whether or not the task is completed.
     */
    updateTask(taskID: number, completed: boolean): Promise<boolean> {
        return this.post<boolean>(`${this.baseURL}/tasks/update`, { taskID: taskID, completed: completed });
    }
}

interface authResult {
    success: boolean,
    token: string,
}