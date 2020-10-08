import React from 'react';
import { BrowserRouter as Router, Switch, Route } from 'react-router-dom';
import './App.css';
import PlayerView from './PlayerView.js';
import TablePage from './TablePage.js';
import TopBar from './TopBar.js';

function App() {
  return (
    <Router>
      <div className='App'>
        <TopBar />
        <div>
            <Switch>
              <Route path='/seat/:id' component={PlayerView} />
              <Route path='/table/:id' component={TablePage} />
              <Route path='/'>
                TODO list of your current drafts
              </Route>
            </Switch>
        </div>
      </div>
    </Router>
  );
}

export default App;
