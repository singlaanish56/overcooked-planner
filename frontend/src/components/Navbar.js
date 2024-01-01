import React from "react"

const Navbar = () => {
    return(
        <nav className="navbar"> 
            <a className="navbar-name" href="/">Overcooked Planner</a>
            <ul>
                <li>
                    <a className="navbar-add-recipe" href="/recipes">Recipes</a>
                </li>
                <li>
                    <a className="navbar-display" href="/display">Display</a>
                </li>
            </ul>
            
        </nav>
    )
}

export default Navbar