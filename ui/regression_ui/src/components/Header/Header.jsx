import './Header.css';

export const Header = ({isRed, children}) => {
    return (
        <div className={`header-container ${isRed ? 'red' : ''}`}>
            <h1 className="header-title rgb-text">{children}</h1>
        </div>
    )
}