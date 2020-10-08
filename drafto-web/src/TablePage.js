import React, { Component } from 'react';
import { Redirect } from 'react-router-dom';
import API from './Api.js';
import { TakeSeatReq } from './service_pb.js';
import TableView from './TableView.js';

class TablePage extends Component {
  constructor(props) {
    super(props);

    this.state = {
      error: null,
    };
  }

  takeSeat() {
    const req = new TakeSeatReq();
    req.setTableId(this.props.match.params.id);
    API.takeSeat(req)
      .then(
        (result) => {
          this.setState({
            data: result,
          });
        },
        (error) => {
          this.setState({
            loaded: true,
            error
          });
        }
      );
  }

  render() {
    if (this.state.data) {
      return <Redirect push to={'/seat/' + this.state.data.seatId} />;
    }
    // TODO: don't render the "take seat" button if the table is full or the player is already at the table
    var err = "";
    if (this.state.error) {
      err = <div>{this.state.error.message}</div>
    }
    return <>{err}<div><button onClick={() => this.takeSeat()}>Take Seat</button></div><TableView id={this.props.match.params.id}/></>;
  }
}

export default TablePage;
