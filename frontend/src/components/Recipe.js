import {useState} from "react"
import DisplayRecipes from "./DisplayRecipes"

const Recipe = () => {
    const [recipes, setRecipes] = useState([])
    const [name, setName] = useState('')
    const [ingredient, setIngredient] = useState('')
    const [ingredients, setIngredients] = useState([])

    const handleNameChange = (event) =>{
        setName(event.target.value)
    }

    const handleIndegdAdd = (event) => {
        setIngredient(event.target.value)
    }

    const handleIndegdinRecipe = (event) => {
        event.preventDefault()
        let found = false
        let y = ingredient.toLowerCase()
        for(let item of ingredients){
            let x = item.toLowerCase()
            if( x == y){
                found = true
                break;
            }
        }
        if(!found)
        {
            setIngredients(ingredients.concat(ingredient))
        }
        
        setIngredient('')
    }

    const handleOnSubmit = (event) => {
        event.preventDefault()
        const recipe = {
            id : name,
            recipeName : name,
            ingredients : ingredients
        }

        setRecipes(recipes.concat(recipe))
        setIngredient('')
        setName('')
        setIngredients([])
    }
    return (
        <div>
        <div className="add-recipe-form">
        <form onSubmit={handleOnSubmit}>
            <div>
                name : <input type="text" value={name} onChange={handleNameChange} />
            </div>
            <div>
                indegredients : <input type="text" value={ingredient} onChange={handleIndegdAdd} />
                <button onClick={handleIndegdinRecipe}>add</button>
            </div>
            <div>
                <button type="submit">Add Recipe</button>
            </div>
        </form>
        </div>
        <DisplayRecipes recipes={recipes} />
        </div>
        

    )
}

export default Recipe