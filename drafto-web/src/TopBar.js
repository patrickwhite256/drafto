import React, { Component } from 'react';
import API from './Api.js';
import { GetCurrentUserReq } from './service_pb.js';

class TopBar extends Component {
  constructor(props) {
    super(props);

    this.state = {
      loggedIn: false,
      username: "",
      avatarURL: "",
    };
  }

  componentDidMount() {
    const req = new GetCurrentUserReq();
    API.getCurrentUser(req)
      .then(
        (result) => {
          this.setState({
            loggedIn: true,
            data: result,
          });
        },
        (error) => {});
  }

  render() {
    var content = <a href="/auth">Login with Discord</a>;
    if (this.state.loggedIn) {
      content = (<><span>{this.state.data.name}</span> <img className="avatar-img" src={this.state.data.avatarUrl} alt="avatar" /></>);
    }

    return <div className="top-bar" >{content}</div>;
  }
}

export default TopBar;
