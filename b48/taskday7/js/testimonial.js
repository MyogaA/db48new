class Parent {
    #deskripsi = ""
    #img = ""

    constructor(deskripsi, img) {
        this.#deskripsi = deskripsi
        this.#img = img
    }

    get deskripsi() {
        return this.#deskripsi
    }

    get img() {
        return this.#img
    }

    get penulis() {
        throw new Error('Harus ada penulis')
    }

    get inHTML() {
        return `
        <div class="grid-item">
            <img src="${this.img}" class=""/>
            <p class="">"${this.deskripsi}"</p>
            <p class="">${this.penulis}</p>
        </div>
        `
    }
}

class Child extends Parent {
    #penulis = ""
    constructor(penulis, deskripsi, img) {
        super(deskripsi, img)
        this.#penulis = penulis
    }

    get penulis() {
        return "penulis : " + this.#penulis
    }
}


const testi1 = new Child("M Yoga Atmajaya", "Sulit Bukan Berarti tidak bisa", "https://images.unsplash.com/photo-1516561318459-78ae3be62459?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHx0b3BpYy1mZWVkfDE2fHhIeFlUTUhMZ09jfHxlbnwwfHx8fHw%3D&auto=format&fit=crop&w=500&q=60")

const testi2 = new Child("Dimas", "Dimas anjay mabar", "https://images.unsplash.com/photo-1688208091969-c916e122eca7?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHx0b3BpYy1mZWVkfDE5fHhIeFlUTUhMZ09jfHxlbnwwfHx8fHw%3D&auto=format&fit=crop&w=500&q=60")

const testi3 = new Child("Stephen", "hayuuuu", "https://images.unsplash.com/photo-1687781834595-2284626b4b95?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHx0b3BpYy1mZWVkfDIyfHhIeFlUTUhMZ09jfHxlbnwwfHx8fHw%3D&auto=format&fit=crop&w=500&q=60")

let data = [testi1, testi2, testi3]
let testiHTML = ""

for (let i = 0; i < data.length; i++) {
    testiHTML += data[i].inHTML
}

document.getElementById("grid-testi").innerHTML = testiHTML
