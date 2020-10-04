import React, { Component } from 'react';
import cardback from './cardback.jpg';

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
      var spanClass ='';
      if (this.props.cardObj.foil) {
        spanClass='foil';
      }
      return <span className={spanClass}><img onClick={this.props.onClick} className='card' src={this.props.cardObj.imageUrl} alt={this.props.cardObj.Name} style={style} /></span>;
    }

    return <span><img className='card' src={cardback} alt='Unrevealed card' style={style} /></span>
  }
}

export default Card;
export const BASE_CARD_HEIGHT = 340;
export const BASE_CARD_WIDTH = 244;
