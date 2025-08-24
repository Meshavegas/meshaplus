export const isEmpty = (str: string) => {   
    // Vérifie si la chaîne est vide, nulle, indéfinie, ne contient que des espaces ou n'est pas une string
    if (typeof str !== 'string') return true
    return str.trim().length === 0
}

export const isNotEmpty = (str: string) => {
    return !isEmpty(str)
}

export const isEmptyObject = (obj: any) => {
    return Object.keys(obj).length === 0
}
