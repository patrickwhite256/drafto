import React, { Component } from 'react';
import API from './Api.js';
import { GetSeatReq, MakeSelectionReq } from './service_pb.js';
import Card from './Card.js';
import CardStack from './CardStack.js';
import TableView from './TableView.js';

class PlayerView extends Component {
  constructor(props) {
    super(props);

    this.state = {
      error: null,
      ws: null,
      loaded: false,
      data: {},
    };
  }

  componentDidMount() {
    this.refreshState();
  }

  startListener(tableID) {
    var socketAddr = 'ws://';
    if (window.location.protocol === "https:") {
      socketAddr = 'wss://';
    }

    socketAddr += window.location.host + '/ws/' + tableID;

    var self = this;
    var ws = new WebSocket(socketAddr);
    ws.onmessage = function(event) {
      if (self.state.error || !self.state.loaded) {
        return
      }

      self.refreshState();
    }

    this.setState({ws: ws});
  }

  componentWillUnmount() {
    if(this.state.ws) {
      this.state.ws.close();
    }
  }

  refreshState() {
    this.setState({loaded: false});

    const req = new GetSeatReq();
    req.setSeatId(this.props.match.params.id);
    API.getSeat(req)
      .then(
        (result) => {
          this.setState({
            loaded: true,
            data: result,
          });

          if(!this.state.ws) {
            this.startListener(result.tableId);
          }
        },
        (error) => {
          this.setState({
            loaded: true,
            error
          });
        }
      );
  }

  selectCard(cardID) {
    const req = new MakeSelectionReq();
    req.setSeatId(this.props.match.params.id);
    req.setCardId(cardID);
    this.setState({loaded: false});
    API.makeSelection(req)
      .then(
        (result) => {
          this.refreshState();
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

    const stackCards = [[],[],[],[],[],[],[]];
    const cards = [];
    for (const cardObj of data.poolList) {
      var stackIdx = 0;
      if (cardObj.coloursList.length > 1) {
        stackIdx = 5;
      } else if (cardObj.coloursList.length === 0) {
        stackIdx = 6;
      } else {
        stackIdx = cardObj.coloursList[0];
      }
      stackCards[stackIdx].push(<Card cardObj={cardObj} revealed={true} />);
      cards.push(<Card cardObj={cardObj} revealed={true} />);
    }

    const poolStacks = [];
    for (var i = 0; i < stackCards.length; i++) {
      if (stackCards[i].length === 0) continue;

      poolStacks.push(<CardStack cards={stackCards[i]} />);
    }

    var pack = <div>No current pack!</div>;
    if (data.currentPack) {
      const cards = [];
      for (const cardObj of data.currentPack.cardsList) {
        cards.push(<Card cardObj={cardObj} revealed={true} onClick={() => this.selectCard(cardObj.id)} />);
      }
      pack = <div>Pack: {cards}</div>;
    }
    return (
      <div>
        <div>
          {pack}
          <div className="cardPool">
            Pool:<br />{poolStacks}
          </div>
        </div>
        Table:
        <TableView id={data.tableId} />
      </div>
    );
  }
}

export default PlayerView;
