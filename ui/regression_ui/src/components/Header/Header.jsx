import './Header.css';

export const Header = ({isRed, children}) => {


    return (
        <div className={isRed ? 'red' : ''}>{children}</div>
    )
}