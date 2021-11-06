import { FunctionalComponent, h } from 'preact';
import Login from '../../components/login';

const Home: FunctionalComponent = () => {
    return (
        <div>
            <h1>Home</h1>
            <Login></Login>
        </div>
    );
};

export default Home;
