import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import API from './Api.js';
import { GetSeatReq } from './service_pb.js';
import { TablePlayer } from './TableView.js';

class SeatView extends Component {
  constructor(props) {
    super(props);

    this.state = {
      loaded: false,
    };
  }

  componentDidMount() {
    const req = new GetSeatReq();
    req.setSeatId(this.props.id);
    API.getSeat(req)
      .then(
        (result) => {
          this.setState({
            loaded: true,
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
    const { error, loaded, data } = this.state;
    if (error) {
      return <div>Error { error.message }</div>;
    } else if (!loaded) {
      return <div>Loading</div>
    }

    console.log(data);
    const seatObj = {
      seatId: data.seat_id,
      packCount: data.packCount,
      poolCount: data.poolList.length,
      poolRevealedCardsList: data.poolList,
      packRevealedCardsList: data.currentPack ? data.currentPack.cardsList : [],
      currentPackCount: data.currentPack ? data.currentPack.cardsList.length : 0,
      playerName: "You",
    };

    return <><Link to={'/seat/' + data.seatId}>Table {data.tableId}</Link><TablePlayer seatObj={seatObj} /></>;
  }
}

export default SeatView;
