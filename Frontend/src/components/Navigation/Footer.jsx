import React from "react";
import { FaInstagram } from "react-icons/fa";
import { FaFacebook } from "react-icons/fa";
import { IoLogoYoutube } from "react-icons/io";
import { HiOutlineMailOpen } from "react-icons/hi";
import logo from "../../assets/Malang.png";

const Footer = () => {
  return (
    <>
      <footer className="bg-gray-200 rounded-lg shadow dark:bg-gray-900 mt-auto">
        {" "}
        <div className="w-full max-w-screen-xl mx-auto">
          <div className="sm:flex sm:items-center sm:justify-between">
            <div className="flex items-center mt-2 mb-2 sm:mb-0 space-x-2 rtl:space-x-reverse">
              {" "}
              <img src={logo} className="h-10 w-auto" alt="Logo" />
              <span className="self-center text-lg font-semibold whitespace-nowrap dark:text-white">
                {" "}
                Desa Ngebruk
              </span>
            </div>
            <ul className="flex flex-wrap items-center mb-6 text-sm font-medium mt-2 text-gray-500 sm:mb-0 dark:text-gray-400 ">
              <li className="md:mr-2 mr-0 ">
                <a
                  href="https://www.instagram.com/desangebruk_poncokusumo?igsh=MWhveHJuZ3VnN29xMg=="
                  className="hover:underline me-4 md:me-6  "
                >
                  <FaInstagram size={27} color=" gray" />
                </a>
              </li>
              <li className="md:mr-2 mr-0">
                <a
                  href="https://web.facebook.com/profile.php?id=61555378417514&_rdc=1&_rdr"
                  className="hover:underline me-4 md:me-6"
                >
                  <FaFacebook size={27} color=" gray" />
                </a>
              </li>
              <li className="md:mr-2 mr-0">
                <a
                  href="https://www.youtube.com/channel/UC9knO2kLi9ooM6qVe_Wz_NA"
                  className="hover:underline me-4 md:me-6"
                >
                  <IoLogoYoutube size={27} color=" gray" />
                </a>
              </li>
              <li className="md:mr-2 mr-0">
                <a
                  href="https://mail.google.com/mail/?view=cm&fs=1&to=ngebrukd@gmail.com&su=Permintaan%20Informasi&body=Halo,%20saya%20ingin%20mebghubungi%20lebih%20lanjut"
                  className="hover:underline me-4 md:me-6"
                >
                  <HiOutlineMailOpen size={27} color=" gray" />
                </a>
              </li>
            </ul>
          </div>
        </div>
      </footer>
    </>
  );
};

export default Footer;
