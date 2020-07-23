function loadDoc(){

    var tablestud = document.querySelector("#tbstud > tbody");
    
    var xhr = new XMLHttpRequest();
    xhr.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 200){
            loadData(tablestud, this.response);
        }
    }; 
    xhr.open("GET", "http://localhost:8080/", true);
    xhr.send();
}

function searchByIndex() {

    var formSearchEl = document.querySelector("#searchForm > input").value;
    if (formSearchEl != 0){
        var tablestud = document.querySelector("#tbstud > tbody");

        var xhr = new XMLHttpRequest();
        xhr.onreadystatechange = function() { 
            if (this.readyState == 4 && this.status == 200){
                loadData(tablestud, this.response);
            }
        };

        xhr.open("GET", "http://localhost:8080/search/" + formSearchEl,true);
        xhr.send();
    } else {
        alert("Enter the index number!");
    }

}

function loadData(tbstud, response){ 
    
    // Removing all current data from table
    while(tbstud.firstChild){
        tbstud.removeChild(tbstud.firstChild);
    }
    
    // Importing new data in table
    jsonobjects = JSON.parse(response);
    jsonobjects.forEach( row => {
         
        const tr = document.createElement("tr");
        const name = document.createElement("td");
        name.textContent = row.name;
        tr.appendChild(name);
        const surname = document.createElement("td");
        surname.textContent = row.surname;
        tr.appendChild(surname);
        const ind = document.createElement("td");
        ind.textContent = row.indexnum;
        tr.appendChild(ind);
        const id = document.createElement("td");
        id.textContent = row.id;
        tr.appendChild(id);        
        
        tbstud.appendChild(tr);
    });
         
}
