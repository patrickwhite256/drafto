import React from 'react';
import { BrowserRouter as Router, Switch, Route, Link } from 'react-router-dom';
import './App.css';
import PlayerView from './PlayerView.js';

function App() {
  return (
    <Router>
      <div className='App'>
        <header className='App-header'>
          
            <Switch>
              <Route path='/player/:id' component={PlayerView} />
              <Route path='/'>
                <Link to='/player/85e70549-2d17-418d-ae8b-8b332df752e6'>whoa</Link>
              </Route>
            </Switch>
        </header>
      </div>
    </Router>
  );
}

export default App;
