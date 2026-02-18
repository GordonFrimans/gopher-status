// src/api/api.js

//gRPC gateway addr

const API_URL = "v1/monitors"; // Адрес твоего gRPC Gateway

//ATTENTION Работает на костыле ... Но работает !
//GET на получение (ничего не принимает)


export const ListMonitors = async () => {
    console.log("GET MONITORS")
    const response = await fetch(`${API_URL}`);
    if (!response.ok) throw new Error("Failed to fetch");
    //УБРАТЬ WARNING
    return response.json();
};


//POST на создание (json принимает)
export const CreateMonitor = async (data) => {
    const response = await fetch(`${API_URL}`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
    });
    if (!response.ok) throw new Error("Failed to create");
    return response.json();
};



//DELETE на удаление по ID (принимает ток id в самом url )
export const DeleteMonitor = async (data) => {
    const response = await fetch(`${API_URL}/${data}`, {
        method: "DELETE",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
    });
    if (!response.ok) throw new Error("Failed to delete");
    return response.json();
};

