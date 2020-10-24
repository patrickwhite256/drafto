import React, { Component } from 'react';
import API from './Api.js';
import { GetCurrentUserReq } from './service_pb.js';
import { BrowserRouter as Router, Switch, Route } from 'react-router-dom';
import './App.css';
import PlayerView from './PlayerView.js';
import TablePage from './TablePage.js';
import TopBar from './TopBar.js';
import SeatView from './SeatView.js';

class DraftList extends Component {
  constructor(props) {
    super(props);

    this.state = {
      loaded: false,
    };
  }

  componentDidMount() {
    const req = new GetCurrentUserReq();
    API.getCurrentUser(req)
      .then(
        (result) => {
          this.setState({
            loaded: true,
            loggedIn: true,
            data: result,
          });
        },
        (error) => {
          this.setState({
            loaded: true,
            loggedIn: false,
          });
        });
  }

  render() {
    const { loggedIn, loaded, data } = this.state;
    if (!loaded) {
      return <div>Loading...</div>;
    }

    if (!loggedIn) {
      return <div>Log in to see your in-progress drafts!</div>;
    }

    const seats = [];
    for (const seatID of data.seatIdsList) {
      seats.push(<><SeatView id={seatID} /></>);
    }

    return <><div>Your Drafts:</div>{seats}</>;
  }
}

function App() {
  return (
    <Router>
      <div className='App'>
        <TopBar />
        <div>
            <Switch>
              <Route path='/seat/:id' component={PlayerView} />
              <Route path='/table/:id' component={TablePage} />
              <Route path='/' component={DraftList} />
            </Switch>
        </div>
      </div>
    </Router>
  );
}

export default App;
