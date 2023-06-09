import { Component, ComponentChild, Fragment, h } from 'preact';
import { TaskState } from './task-state';
import style from './style.css';
import { APIService } from '../../services/api.service';
import { SessionService } from '../../services/session.service';
import { route } from 'preact-router';

class TaskList extends Component<{}, TaskState> {
    api: APIService;

    render(): ComponentChild {
        let results;
        if (this.state.tasks.length === 0) {
            results = <p>No tasks</p>;
        } else {
            results = <div>
                {this.state.tasks.map(item => (
                    <label class={style.container} key={item.id}>
                        {item.title}
                        <input type="checkbox" checked={item.completedAt != null} onClick={(e) => this.checkChanged(e, item.id)}></input>
                        <span class={style.checkmark}></span>
                    </label>
                ))}
            </div>;
        }
        return (
            <div>
                <h1>Tasks</h1>
                <button onClick={this.downloadTasks}>Download Tasks</button>
                {!this.state.initialLoad ? <div class={style.loading}></div> : results}
                <form onSubmit={this.submitTask}>
                    <input type="text" placeholder="New task" value={this.state.newTaskName} onInput={this.taskNameChanged} autoFocus={true}></input>
                    <button type="submit" disabled={this.state.newTaskName.length === 0}>Add Task</button>
                </form>
            </div>
        );
    }

    constructor() {
        super();
        this.api = new APIService();
        this.state = {
            initialLoad: false,
            tasks: [],
            newTaskName: '',
        };
        if (!SessionService.isAuthenticated()) {
            route('/', true);
        }
    }

    componentDidMount() {
        // Grab tasks from API
        this.refreshTaskList();
    }

    refreshTaskList() {
        this.api.getTaskList().then(tasks => {
            this.setState({
                initialLoad: true,
                tasks: tasks,
            })
        });
    }

    submitTask = (e: Event) => {
        e.preventDefault();
        this.api.addTask(this.state.newTaskName).then(newID => {
            this.setState({ newTaskName: '' });
            this.refreshTaskList();
        });
    }

    taskNameChanged = (e: any) => {
        const value = e.target.value;
        this.setState({ newTaskName: value })
    }

    checkChanged = (e: any, itemID: number) => {
        this.api.updateTask(itemID, e.target.checked).then(() => {
            // Give it time just so user sees check thange, then sees task move.
            setTimeout(() => {
                this.refreshTaskList()
            }, 200);
        });
    }

    downloadTasks = () => {
        this.api.getTaskListExcel().then(res => {
            const url = window.URL.createObjectURL(res);
            const link = document.createElement('a');
            link.href = url;
            link.setAttribute('download', 'tasks.xlsx');
            document.body.appendChild(link);
            link.click();
        }).catch(err => {
            console.log(err);
        });
        /*
        
        const url = window.URL.createObjectURL(new Blob([res.data]));
        const link = document.createElement('a');
        link.href = url;
        link.setAttribute('download', 'tasks.xlsx');
        document.body.appendChild(link);
        link.click();
        */
    }
}

export default TaskList;