import { Component, ComponentChild, h } from 'preact';
import { route } from 'preact-router';
import { APIService } from '../../services/api.service';

interface loginState {
    email: string;
    password: string;
}

class Login extends Component<{}, loginState> {
    api: APIService;

    render(): ComponentChild {
        return (
            <form onSubmit={this.logIn}>
                <label for="email">Email:</label>
                <input name="email" type="email" placeholder="Email" value={this.state.email} onInput={this.handleInputChange} autoFocus={true}></input>
                <label for="password">Password:</label>
                <input name="password" type="password" placeholder="Password" value={this.state.password} onInput={this.handleInputChange}></input>
                <button type="submit" disabled={this.state.email.length == 0 || this.state.password.length ===0}>Log In</button>
                <button type="button" onClick={this.newAccount}>Create Account</button>
            </form>
        );
    }

    constructor() {
        super();
        this.api = new APIService();
        this.state = {
            email: '',
            password: '',
        };
    }
    
    handleInputChange= (event: any) => {
        const target = event.target;
        const value = target.value;
        const name = target.name;
        this.setState({
            [name]: value
        });
    }

    logIn = (e: Event) => {
        e.preventDefault();
        this.api.logIn(this.state.email, this.state.password).then(result => {
            if (result === true) {
                route('/tasks', true);
            }
        });
        
        //window.location.href = '/tasks';
    }

    newAccount = (e: Event) => {
        e.preventDefault();
        this.api.createAccount(this.state.email, this.state.password).then(result => {
            if (result === true) {
                route('/tasks', true);
            }
        });
    }
}

export default Login;