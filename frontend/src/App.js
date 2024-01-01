import React from "react"
import Navbar from "./components/Navbar"
import {Routes, Route} from "react-router-dom"
import Recipe from "./components/Recipe"
import Home from "./components/Home"
import Display from "./components/Display"
const App = () => {
    
    return (
        <div>
        <Navbar />
        <div className="container">
         <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/recipes" element={<Recipe />} />
            <Route path="/display" element={<Display />} />
        </Routes>           
        </div>
        </div>




       
    )
}

export default App