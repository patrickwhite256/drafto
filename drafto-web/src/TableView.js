import React, { Component } from 'react';
import API from './Api.js';
import { GetDraftStatusReq } from './service_pb.js';
import Card from './Card.js';
import CardStack from './CardStack.js';

class TableView extends Component {
  constructor(props) {
    super(props);

    this.state = {
      error: null,
      loaded: false,
      data: {},
    };
  }

  componentDidMount() {
    this.refreshState();
  }

  refreshState() {
    const req = new GetDraftStatusReq();
    req.setTableId(this.props.id);
    API.getDraftStatus(req)
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

    const nRows = Math.ceil(data.seatsList.length / 2);
    const rows = [];
    for (var i = 0; i < nRows; i++) {
      const playerDivs = [<div style={{'float': 'left'}}><TablePlayer seatObj={data.seatsList[2*i]} /></div>];
      if (data.seatsList.length > 2*i + 1) {
        playerDivs.push(<div style={{'float': 'right'}}><TablePlayer seatObj={data.seatsList[2*i+1]} /></div>);
      }
      rows.push(<div>{playerDivs}<div style={{clear:'both'}} /></div>);
    }

    return (
      <div>{rows}</div>
    );
  }
}

class TablePlayer extends Component {
  render() {
    const nUnrevealed = this.props.seatObj.poolCount - this.props.seatObj.poolRevealedCardsList.length;
    const poolCards = [];
    for (const cardObj of this.props.seatObj.poolRevealedCardsList) {
      poolCards.push(<Card cardObj={cardObj} revealed={true} scale={0.5} />);
    }
    for (var i = 0; i < nUnrevealed; i++) {
      poolCards.push(<Card revealed={false} scale={0.5}/>);
    }

    const stacks = [];
    for (i = 0; i < Math.ceil(poolCards.length / 10); i++) {
      const stackCards = [];
      for (var j = 10 * i; j < poolCards.length && j < 10 * (i+1); j++) {
        stackCards.push(poolCards[j]);
      }
      stacks.push(<CardStack cards={stackCards} scale={0.5} />);
    }

    const nPackUnrevealed = this.props.seatObj.currentPackCount - this.props.seatObj.packRevealedCardsList.length;
    const packCards = [];
    for (const cardObj of this.props.seatObj.packRevealedCardsList) {
      packCards.push(<Card cardObj={cardObj} revealed={true} scale={0.5} />);
    }
    for (i = 0; i < nPackUnrevealed; i++) {
      packCards.push(<Card revealed={false} scale={0.5}/>);
    }

    return (
      <div className="player-box">
        <div>{this.props.seatObj.playerName} - Current Packs: {this.props.seatObj.packCount}</div>
        <span className="player-pool">
          Pool:<br/>
          {stacks}
        </span>
        <span className="player-pack">
          Pack:<br/>
          <CardStack scale={0.5} horizontal={true} cards={packCards} />
        </span>
      </div>
    );
  }
}

export default TableView;
export { TablePlayer };
