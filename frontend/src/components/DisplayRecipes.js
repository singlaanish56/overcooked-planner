import React from "react"

const DisplayRecipe = ({key,recipe}) => {
    return (
        <div>
            {recipe.recipeName}
            <ul>
            {recipe.ingredients.map(ingredient =>
                <li key="ingredient">{ingredient}</li>
            )

            }
            </ul>

        </div>
    )
}

const DisplayRecipes = ({recipes}) =>{
    return (
        <div>
            {recipes.map(recipe =>
                <DisplayRecipe key={recipe.id} recipe={recipe} />
            )}
        </div>
    )
}
export default DisplayRecipes