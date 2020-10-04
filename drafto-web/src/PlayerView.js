import React, { Component } from 'react';
import API from './Api.js';
import { GetSeatReq, MakeSelectionReq } from './service_pb.js';
import Card from './Card.js';
import CardStack from './CardStack.js';
import TableView from './TableView.js';

class PlayerView extends Component {
  constructor(props) {
    super(props);

    this.id = props.match.params.id;
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
    const req = new GetSeatReq();
    req.setSeatId(this.id);
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

  selectCard(cardID) {
    const req = new MakeSelectionReq();
    req.setSeatId(this.id);
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

    const cards = [];
    for (const cardObj of data.poolList) {
      cards.push(<Card cardObj={cardObj} revealed={true} />);
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
            Pool:<br /><CardStack cards={cards} />
          </div>
        </div>
        <TableView id={data.tableId} />
      </div>
    );
  }
}

export default PlayerView;
