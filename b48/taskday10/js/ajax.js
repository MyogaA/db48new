function allcard (){
    let inHtml = ""

    cardData.forEach((element)=> {
        inHtml += ` <div class="grid-item">
        <img src="${element.img}" class=""/>
        <p class="">"${element.deskripsi}"</p>
        <p class="">${element.penulis}</p>
        <p class="">${element.rating}<i class="fa-solid fa-star"></i></p>
    </div>`
    })
    document.getElementById('grid-testi').innerHTML = inHtml
}

function filterRating (rating){
    let filterInhtml = ""

    const filterdata = cardData.filter((element)=>{
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
const promise = new Promise((berhasil, gagal) => {
    const xhr = new XMLHttpRequest()

    xhr.open("GET", "https://api.npoint.io/7daf4b796facc7155627", true)
    xhr.onload = function () {
        if (xhr.status === 200) {
            berhasil(JSON.parse(xhr.responseText))
        } else if (xhr.status >= 400) {
            gagal("Error loading data")
        }
    }
    xhr.onerror = function () {
        gagal("Network error")
    }
    xhr.send()
})
let cardData = []

async function getData() {
    try {
        const response = await promise
        console.log(response)
        cardData = response
        allcard()
    } catch (err) {
        console.log(er)
    }
}

getData()