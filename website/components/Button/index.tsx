import React, { ReactNode } from 'react';
import classnames from 'classnames';

export enum ButtonType {
  secondary = 'secondary',
  large = 'large',
}

interface ButtonProps {
  children: ReactNode;
  type?: ButtonType;
  onClick(): void;
}

const Button: React.FC<ButtonProps> = ({type, children, onClick}) => {
  const buttonClasses = classnames(`button ${type && `button--${type}`}`);
  return (
    <button type="button" className={buttonClasses} onClick={onClick}>
      {children}
    </button>
  )
}

export default Button;
