import React, { useState } from "react";

const ProfileDesa = () => {
  const [showSejarah, setShowSejarah] = useState(false);
  const [showVisiMisi, setShowVisiMisi] = useState(false);
  const [showStruktur, setShowStruktur] = useState(false);

  return (
    <div className="flex flex-col items-center min-h-screen py-10 bg-gray-50 mt-16">
      <div className="container mx-auto p-4 max-w-full w-3/4">
        <div className="mb-4">
          <button
            onClick={() => setShowSejarah(!showSejarah)}
            className="w-full text-blue-500 hover:text-white border-2 border-blue-500 hover:bg-blue-700 font-medium py-2 px-4 rounded focus:outline-none focus:shadow-outline"
          >
            Sejarah Desa
          </button>
          {showSejarah && (
            <div className="mt-2 p-4 border border-blue-200 bg-orange-50 rounded shadow-xl shadow-slate-300">
              <p className="indent-4">
                Sejarah Desa Ngebruk berasal dari sebuah cerita : Pada jaman
                dahulu ada seorang pengembara, bernama Mbah Djiman yang berasal
                dari Daerah Mataram dan Mbah Salmah yang berasal dari suku Jawa
                mereka adalah pasangan suami istri pada masa itu Desa Ngebruk
                masih hutan belantara sejak saat itu mereka mulai membuat
                perkampungan dan mereka bersepakat bahwa kelak kalau ada
                ramainya jaman, daerah ini dinamakan Desa Ngebruk yang berasal
                dari tanah Brugkan sebelah makam Islam. Mula-mula Desa Ngebruk
                cuma ada satu rumah di sebelah timur dengan sebutan kampung
                nongkopait yang berasal dari pohon nangka dan konon katanya
                nangkanya pahit, dengan berjalanya waktu lalu terbentuklah Desa,
                terdiri dari 5 RW, 55 RT. Bahasa sehari-hari yang digunakan oleh
                warga desa adalah bahasa jawa. Mereka penduduk desa menganut
                agama islam hingga sekarang. Adapun kepala desa yang pernah
                menjabat hingga sekarang adalah sebagai berikut:
              </p>
              <ol className="list-decimal p-4">
                <li>Aris Noto: 1916 - 1940</li>
                <li>Karto Joyo: 1940 - 1959</li>
                <li>Suriyono: 1959 - 1973</li>
                <li>Lasim: 1973-1976 (Pemilihan)</li>
                <li>Ruba'i: 1976-1987 (Pemilihan)</li>
                <li>Sayutik: 1987-2007 (Pemilihan)</li>
                <li>Sodikin: 2007-2013 (Pemilihan)</li>
                <li>Pujiono: 2013-2025 (Pemilihan)</li>
              </ol>
              <p className="italic">
                (Sumber data: Data Profil Desa Ngebruk 2019)
              </p>
            </div>
          )}
        </div>

        {/* Visi dan Misi */}
        <div className="mb-4">
          <button
            onClick={() => setShowVisiMisi(!showVisiMisi)}
            className="w-full text-blue-500 hover:text-white border-2 border-blue-500 hover:bg-blue-700 font-medium py-2 px-4 rounded focus:outline-none focus:shadow-outline"
          >
            Visi dan Misi
          </button>
          {showVisiMisi && (
            <div className="mt-2 p-4 border border-blue-200 bg-orange-50 rounded shadow-xl shadow-slate-300">
              <h1 className=" font-bold text-2xl text-center">Visi</h1>
              <p className="indent-4">
                Visi merupakan pandangan jauh ke depan, ke mana dan bagaimana
                Desa Ngebruk harus dibawa dan berkarya agar konsisten dan dapat
                eksis, antisipatif, inovatif serta produktif. Didalam
                meningkatkan kesejahteraan masyarakat, Desa Ngebruk memiliki
                visi yang dirumuskan berdasarkan potensi yang dimiliki oleh
                masyarakat Desa Ngebruk dan implementasiannya dituangkan dalam
                Misi yang merupakan rumusan langkah-langkah pencapaiannya.
              </p>
              <p className="indent-4">
                Visi adalah gambaran mengenai masa depan dan masa sekarang
                dengan dasar logika dan makna secara bersama. Selanjutnya
                memberi ilham dan naluri yang mensyaratkan harapan dan
                kebanggaan apabila berhasil. Untuk itulah pemerintah Desa
                Ngebruk dalam mencapai cita- citanya memiliki Visi: “Terwujudnya
                Masyarakat adil, makmur, aman, tentram dan sejahtera melalui
                peningkatan kualitas sumber daya manusia yang terdidik dengan
                didukung pengembangan ekonomi berbasis sumber daya alam dan
                dilandasi kehidupan beragama yang kuat”
              </p>
              <p className="indent-4">
                Melalui visi ini diharapkan masyarakat menemukan gambaran
                kondisi masa depan yang lebih baik (ideal) dan merupakan potret
                keadaan yang ingin dicapai, dibanding dengan kondisi yang ada
                saat ini. Melalui rumusan visi ini diharapkan mampu memberikan
                arah perubahan masyarakat pada keadaan yang lebih baik,
                menumbuhkan kesadaran masyarakat untuk mengendalikan dan
                mengontrol perubahan-perubahan yang akan terjadi, mendorong
                masyarakat untuk meningkatkan kinerja yang lebih baik,
                menumbuhkan kompetisi sehat pada anggota masyarakat, menciptakan
                daya dorong untuk perubahan serta mempersatukan anggota
                masyarakat.
              </p>
              <h1 className=" font-bold text-2xl text-center">Misi</h1>
              <p className="indent-4">
                Misi adalah rumusan umum mengenai upaya-upaya yang akan
                dilaksanakan untuk mewujudkan visi. Misi berfungsi sebagai
                pemersatu gerak, langkah dan tindakan nyata bagi segenap
                komponen penyelenggara pemerintahan tanpa mengabaikan mandat
                yang diberikannya.
              </p>
              <p className="indent-4">
                Hakekat misi merupakan turunan dari visi yang akan menunjang
                keberhasilan tercapainya sebuah visi. Dengan kata lain Misi
                merupakan penjabaran lebih operatif dari Visi. Penjabaran dari
                visi ini diharapkan dapat mengikuti dan mengantisipasi setiap
                terjadinya perubahan situasi dan kondisi lingkungan di masa yang
                akan datang dari usaha-usaha mencapai Visi desa selama masa 5
                (lima) tahun.
              </p>
              <p className="indent-4">
                Untuk meraih Visi desa seperti yang sudah dijabarkan di atas,
                dengan mempertimbangan potensi dan hambatan baik internal maupun
                eksternal, maka disusunlah Misi desa sebagai berikut:
              </p>
              <ol className="list-decimal list-inside space-y-2">
                <li>
                  1. Melaksanakan/mengamalkan ajaran agama dalam kehidupan
                  bermasyarakat berbangsa dan bernegara sebagai wujud
                  peningkatan keimanan dan ketaqwaan kepada Tuhan Yang Maha Esa.
                </li>
                <li>
                  Mewujudkan dan mendorong terjadinya usaha-usaha kerukunan
                  antar dan intern warga masyarakat yang disebabkan karena
                  adanya perbedaan agama, keyakinan, organisasi, dan lainnya
                  dalam suasana saling menghargai dan menghormati.
                </li>
                <li>
                  Mengembangkan kehidupan masyarakat untuk terwujudnya tatanan
                  masyarakat yang taat kepada peraturan perundang-undangan dalam
                  rangka meningkatkan kehidupan masyarakat yang aman, tertib,
                  tentram dan damai serta meningkatakan persatuan dan kesatuan
                  dalam wadah negara kesatuan Republik Indonesia.
                </li>

                <li>
                  Membangun dan meningkatkan hasil pertanian dengan jalan
                  penataan pengairan, perbaikan jalan sawah / jalan usaha tani,
                  pemupukan, dan polatanam yang baik.
                </li>
                <li>
                  Pengembangan sektor pertanian dan perdagangan yang
                  berorientasi pada mekanisme pasar.
                </li>
                <li>Menumbuhkembangkan usaha kecil dan menengah.</li>
                <li>
                  Pemberdayaan ekonomi masyarakat khususnya UMKM (Usaha Kecil
                  Menengah dan Mikro) yang berdaya saing tinggi.
                </li>
                <li>
                  Membangun dan mendorong usaha-usaha untuk pengembangan dan
                  optimalisasi sektor pertanian, perkebunan, peternakan, dan
                  perikanan, baik tahap produksi maupun tahap pengolahan
                  hasilnya
                </li>
                <li>
                  Meningkatkan kemajuan dan kemandirian melalui penyelenggaraan
                  otonomi desa yang bertanggung jawab dan didukung dengan yang
                  bersih, transparan penyelenggaran profesional pemerintahan
                </li>
              </ol>
            </div>
          )}
        </div>

        {/* Struktur Organisasi */}
        <div className="mb-4">
          <button
            onClick={() => setShowStruktur(!showStruktur)}
            className="w-full text-blue-500 hover:text-white border-2 border-blue-500 hover:bg-blue-700 font-medium py-2 px-4 rounded focus:outline-none focus:shadow-outline"
          >
            Struktur Organisasi
          </button>
          {showStruktur && (
            <div className="mt-2 p-4 border border-blue-200 bg-orange-50 rounded shadow-xl shadow-slate-300">
              <ol className="list-decimal list-inside space-y-2">
                <li>Kepala Desa: Sanam</li>
                <li>Sekretaris Desa: Tomi Aris Wibowo</li>
                <li>Kepala Urusan Tata Usaha dan Umum: Rudyantoko</li>
                <li>Kepala Urusan Perencanaan: Wari</li>
                <li>Kepala Urusan Keuangan: M. Taufik</li>
                <li>Kepala Seksi Pelayanan: Singgik</li>
                <li>Kepala Seksi Pemerintahan: Sutris</li>
                <li>Kepala Seksi Kesejahteraan: Muhaimin</li>
                <li>Kepala Dusun 1: Arif</li>
                <li>Kepala Dusun 2: Supriono</li>
              </ol>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default ProfileDesa;
