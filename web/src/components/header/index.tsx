import { FunctionalComponent, h } from 'preact';
import { Link } from 'preact-router/match';
import { SessionService } from '../../services/session.service';
import style from './style.css';

const Header: FunctionalComponent = () => {
    return (
        <header class={style.header}>
            <div class="maxWidth">
                <h1>Task App</h1>
                <nav>
                    <Link activeClassName={style.active} href="/">Home</Link>
                    <Link activeClassName={style.active} href="/tasks">Tasks</Link>
                </nav>
            </div>
        </header>
    );
};

export default Header;
