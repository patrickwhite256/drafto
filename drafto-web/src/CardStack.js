import React, { Component } from 'react';
import { BASE_CARD_HEIGHT, BASE_CARD_WIDTH} from './Card.js';

class CardStack extends Component {
  constructor(props) {
    super(props);

    this.scale = 1;
    if (props.scale) {
      this.scale = props.scale;
    }
  }

  render() {
    const cards = [];
    for (var i = 0; i < this.props.cards.length; i++) {
      const cardContainer = <CardContainer scale={this.scale} horizontal={this.props.horizontal} cardN={i} card={this.props.cards[i]} />
      cards.push(cardContainer);
    }

    var style = {
      'width': BASE_CARD_WIDTH*this.scale,
      'height': (BASE_CARD_HEIGHT+40*(cards.length-1))*this.scale,
    };
    if (this.props.horizontal) {
      style = {
        'width': (BASE_CARD_WIDTH+40*(cards.length-1))*this.scale,
        'height': BASE_CARD_HEIGHT*this.scale,
      };
    }

    return (
      <span className='card-stack' style={style}>
       {cards}
      </span>
    );
  }
}

class CardContainer extends Component {
  constructor(props) {
    super(props)
   
    this.state = {z: props.cardN};
    this.scale = 1;
    if (props.scale) {
      this.scale = props.scale;
    }
  }

  bringToFront() {
    this.setState({z: 9999});
  }

  resetZ() {
    this.setState({z: this.props.cardN});
  }

  render() {
    var style = {
      'top': 40*this.scale*this.props.cardN,
      'zIndex': this.state.z,
    };
    if (this.props.horizontal) {
      style = {
        'left': 40*this.scale*this.props.cardN,
        'zIndex': this.state.z,
      };
    }
    return (
      <span style={style} onMouseOver={()=>this.bringToFront()} onMouseLeave={()=>this.resetZ()}>
        {this.props.card}
      </span>
    );
  }
}

export default CardStack;
