import React from "react"

const Navbar = () => {
    return(
        <div className="navbar"> 
            <a className="navbar-name" href="#name">Overcooked Planner</a>
            <a className="navbar-add-recipe" href="#recipe">Add Recipes</a>
            <a className="navbar-display" href="#display">Display</a>
        </div>
    )
}

export default Navbar