import React, { Component } from 'react';

const CARD_BACK_URL = 'https://c1.scryfall.com/file/scryfall-card-backs/normal/59/597b79b3-7d77-4261-871a-60dd17403388.jpg';

class Card extends Component {
  constructor(props) {
    super(props);

    this.scale = 1;
    if (props.scale) {
      this.scale = props.scale;
    }
  }

  render() {
    const style = {'height': BASE_CARD_HEIGHT * this.scale, 'width': BASE_CARD_WIDTH * this.scale};

    if (this.props.revealed) {
      return <img onClick={this.props.onClick} className='card' src={this.props.cardObj.imageUrl} alt={this.props.cardObj.Name} style={style} />
    }

    return <img className='card' src={CARD_BACK_URL} alt='Unrevealed card' style={style} />
  }
}

export default Card;
export const BASE_CARD_HEIGHT = 340;
export const BASE_CARD_WIDTH = 244;
