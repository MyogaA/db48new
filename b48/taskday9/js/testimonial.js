function allcard (){
    let inHtml = ""

    Data.forEach((element)=> {
        inHtml += ` <div class="grid-item">
        <img src="${element.img}" class=""/>
        <p class="">"${element.deskripsi}"</p>
        <p class="">${element.penulis}</p>
        <p class="">${element.rating}<i class="fa-solid fa-star"></i></p>
    </div>`
    })
    document.getElementById('grid-testi').innerHTML = inHtml
}
const Data = [
    {
        penulis: "M Yoga Atmajaya",
        deskripsi: "Sulit bukan berarti tidak bisa",
        img: "https://images.unsplash.com/photo-1516561318459-78ae3be62459?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHx0b3BpYy1mZWVkfDE2fHhIeFlUTUhMZ09jfHxlbnwwfHx8fHw%3D&auto=format&fit=crop&w=500&q=60",
        rating: 5
    },
    {
        penulis: "Dimas",
        deskripsi: "Dimas anjay mabar",
        img: "https://images.unsplash.com/photo-1688208091969-c916e122eca7?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHx0b3BpYy1mZWVkfDE5fHhIeFlUTUhMZ09jfHxlbnwwfHx8fHw%3D&auto=format&fit=crop&w=500&q=60",
        rating: 3
    },
    {
        penulis: "Stephen",
        deskripsi: "hayuuuuu",
        img: "https://images.unsplash.com/photo-1687781834595-2284626b4b95?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHx0b3BpYy1mZWVkfDIyfHhIeFlUTUhMZ09jfHxlbnwwfHx8fHw%3D&auto=format&fit=crop&w=500&q=60",
        rating: 2
    },
    {
        penulis: "Basid",
        deskripsi: "hiyahiyahiya",
        img: "https://images.unsplash.com/photo-1688208091969-c916e122eca7?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHx0b3BpYy1mZWVkfDE5fHhIeFlUTUhMZ09jfHxlbnwwfHx8fHw%3D&auto=format&fit=crop&w=500&q=60",
        rating: 2
    },
    {
        penulis: "Dika",
        deskripsi: "Gokil",
        img: "https://images.unsplash.com/photo-1687781834595-2284626b4b95?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHx0b3BpYy1mZWVkfDIyfHhIeFlUTUhMZ09jfHxlbnwwfHx8fHw%3D&auto=format&fit=crop&w=500&q=60",
        rating: 1
    },
]
allcard()

function filterRating (rating){
    let filterInhtml = ""

    const filterdata = Data.filter((element)=>{
        return element.rating === rating
    })
    filterdata.forEach((element)=> {
        filterInhtml += `<div class="grid-item">
        <img src="${element.img}" class=""/>
        <p class="">"${element.deskripsi}"</p>
        <p class="">${element.penulis}</p>
        <p class="">${element.rating}<i class="fa-solid fa-star"></i></p>
    </div>`
    })
    document.getElementById("grid-testi").innerHTML = filterInhtml
}