import { FunctionalComponent, h } from 'preact';
import { Route, Router } from 'preact-router';

import Home from './home';
import NotFoundPage from './notfound';
import Header from '../components/header';
import TaskList from './task-list';

const App: FunctionalComponent = () => {
    return (
        <div id="preact_root">
            <Header />
            <div class="content-wrap">
                <Router>
                    <Route path="/" component={Home} />
                    <Route path="/tasks" component={TaskList} />
                    <NotFoundPage default />
                </Router>
            </div>
        </div>
    );
};

export default App;
