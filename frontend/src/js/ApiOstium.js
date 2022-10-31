import ApiBase from "./ApiBase";

// Method to get our current login details from the backend
export const getLogin = async () => {
    try {
        const response = await ApiBase.get("/login");
        return response;
    } catch (error) {
        console.error(error);
    }
};

// Method to log out from the backend
export const deleteLogin = async () => {
    try {
        const response = await ApiBase.delete("/login");
        return response;
    } catch (error) {
        console.error(error);
    }
};

// Method to log into the backend
export const postLogin = async (username, password) => {
    try {
        let request = {
            username: username,
            password: password
        }
        const response = await ApiBase.post("/login", JSON.stringify(request));
        return response;
    } catch (error) {
        console.error(error);
    }
};

// Method to load a set of available objects of a type
export const getObjects = async (collection) => {
    try {
        const response = await ApiBase.get("/" + collection);
        return response;
    } catch (error) {
        console.error(error);
    }
};
