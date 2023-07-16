//  kode ini mengambil referensi elemen-elemen HTML yang akan digunakan,
//  seperti form, input, dan elemen kontainer untuk menampilkan data proyek.
const projectForm = document.getElementById('project-form');
const projectNameInput = document.getElementById('projectname');
const startDateInput = document.getElementById('startdate');
const endDateInput = document.getElementById('enddate');
const descriptionInput = document.getElementById('message');
const imageInput = document.getElementById('input-img');
const projectCards = document.getElementById('project-card');


// Array untuk menyimpan data project
const projects = [];

// Event listener saat form disubmit
// sebuah event listener ditambahkan pada form untuk menangani saat form disubmit.
// Event listener ini mencegah pengiriman form secara default menggunakan 
projectForm.addEventListener('submit', function (event) {
  event.preventDefault(); // Mencegah form dari pengiriman default

  // Mengambil nilai dari form
  // Saat form disubmit, nilai-nilai dari input form diambil menggunakan value properti dari elemen input.
  const projectName = projectNameInput.value;
  const startDate = startDateInput.value;
  const endDate = endDateInput.value;
  const description = descriptionInput.value;
  const imageFile = imageInput.files[0];

  // untuk membaca berkas gambar dan mengubahnya menjadi URL
  const reader = new FileReader(); // Untuk membaca dan mengubah gambar menjadi URL, FileReader digunakan.
  reader.onload = function () { // metode onload akan dipanggil saat pembacaan selesai.
    const imageUrl = reader.result; // Di dalam onload, URL gambar dibaca dari reader.result.

    // Mengambil nilai dan data dari checkbox
    const nodejs = document.getElementById('nodejs').checked ? '<i class="fa-brands fa-node-js"></i>' : '';
    const Reactjs = document.getElementById('reactjs').checked ? '<i class="fa-brands fa-react"></i>' : '';
    const nextjs = document.getElementById('nextjs').checked ? '<i class="fa-brands fa-android"></i>' : '';
    const typescript = document.getElementById('typescript').checked ? '<i class="fa-brands fa-java"></i>' : '';

    // Membuat objek data project baru
    // Sebuah objek data proyek baru dibuat dengan nilai-nilai yang diambil dari form.
    const newProject = {
      projectName: projectName,
      startDate: startDate,
      endDate: endDate,
      description: description,
      imageUrl: imageUrl,
      technologies: [nodejs, Reactjs, nextjs, typescript].filter(icon => icon !== '')
    };

    // Menambahkan data project baru ke dalam array projects
    // Objek data proyek baru ditambahkan ke dalam array projects.
    projects.push(newProject);

    // Menampilkan semua data proyek menggunakan innerHTML dan loop
    // Data proyek ditampilkan menggunakan innerHTML dan loop.
    //  Sebuah string HTML dibangun untuk setiap proyek dalam array projects,
    //  termasuk gambar, nama proyek, detail proyek, deskripsi, dan tombol aksi.
    projectCards.innerHTML = '';
    for (let i = 0; i < projects.length; i++) {
      const project = projects[i];

      // Menghitung durasi waktu dalam hari
      const startDateObj = new Date(project.startDate);
      const endDateObj = new Date(project.endDate);
      const timeDiff = Math.abs(endDateObj - startDateObj);
      const durationDays = Math.ceil(timeDiff / (1000 * 60 * 60 * 24));

      // Membuat elemen HTML untuk durasi waktu dan waktu hari
      const durationHTML = `<span class="duration">Durasi: ${durationDays} hari</span>`;
      const timeHTML = `<span class="time">${project.startDate} - ${project.endDate}</span>`;

      projectCards.innerHTML += `
        <div class="card-item">
          <div class="project-image">
            <img src="${project.imageUrl}" alt="">
          </div>
          <div class="text-content">
            <h4 style="float: left; padding-left: 20px;"><a style="text-decoration: none; color: black;" href="detail.html">${project.projectName}</a></h4>
            <div class="project-details">
              ${day} / ${bulan[month]} / ${year} ${hour}:${minute}
              ${durationHTML}
              ${timeHTML}
            </div>
            <p>${project.description}</p>
          </div>
          <div class="logo-bahasa">
            ${project.technologies.join('')}
          </div>
          <div class="btn-project">
            <button type="button" class="btn1">edit</button>
            <button type="button" class="btn2">delete</button>
          </div>
        </div>
      `;
    }

    // Mengosongkan ulang nilai input form setelah data ditampilkan
    // Setelah data ditampilkan, nilai-nilai input form dikosongkan.
    projectNameInput.value = '';
    startDateInput.value = '';
    endDateInput.value = '';
    descriptionInput.value = '';
    imageInput.value = '';

    // Mengosongkan URL objek saat gambar tidak lagi diperlukan
    // URL objek gambar dihapus menggunakan URL.revokeObjectURL(imageUrl) untuk menghemat memori.
    URL.revokeObjectURL(imageUrl);
  };
  reader.readAsDataURL(imageFile);
  
  // kode yang mengatur tampilan tanggal saat ini dalam format tanggal, bulan, tahun, jam, dan menit.
  //  Nilai-nilai ini digunakan untuk menampilkan tanggal saat ini pada elemen HTML.
  // Kode tersebut juga mencakup loop dan pemrosesan array bulan untuk menampilkan nama bulan dalam bahasa Indonesia.
  bulan = ["Januari", "Februari", "Maret", "April", "Mei", "Juni", "Juli", "Agustus", "September", "Oktober", "November", "Desember"];

  let date = new Date();
  let day = date.getDate();
  let month = date.getMonth();
  let year = date.getFullYear();
  let hour = date.getHours();
  let minute = date.getMinutes();

  if (day < 10) day = `0${day}`;
  if (hour < 10) hour = `0${hour}`;
  if (minute < 10) minute = `0${minute}`;


  bulan.forEach((value, index) => {
    console.log(value, index);
  });
  for (let i = 0; i < bulan.length; i++) {  
    console.log()
  }
});
projectCards.addEventListener('click', function (event) {
  if (event.target.classList.contains('btn2')) {
    const card = event.target.closest('.card-item'); // Mendapatkan elemen card terdekat yang mengandung tombol delete
    const index = Array.from(projectCards.children).indexOf(card); // Mendapatkan indeks elemen card dalam projectCards

    if (index !== -1) {
      projects.splice(index, 1); // Menghapus data proyek dari array projects berdasarkan indeksnya
      card.remove(); // Menghapus elemen card dari tampilan Html
    }
  }
});
