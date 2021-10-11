import logo from './logo.svg';
import './App.css';
import {BrowserRouter as Router, Route, Link} from 'react-router-dom';
import Protected from './Protected';

function App() {
  return (
    <div className="App">
        <Router>
          <header className="App-header">
            <Route exact path='/'>
              <img src={logo} className="App-logo" alt="logo" />
              <p>
                Hello Anonymous. Please login <Link to='/protected'>here</Link>
              </p>
            </Route>
            <Route path='/protected'>
                <Protected />
            </Route>
          </header>
        </Router>
    </div>
  );
}

export default App;
