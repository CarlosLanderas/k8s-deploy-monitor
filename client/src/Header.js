import React from "react";

import k8slarge from "./k8s-large.png"


const Header = () => {
    return (
    <header className="header">
    <div className="container">
        <img className="logo" src={k8slarge}/>
        <div style={{marginLeft: "6px"}}>deployment monitor</div>
    </div>
        
    </header>
    )   
 
}

export default Header;
  
  