"use client";
import dynamic from "next/dynamic";
import React, { useContext } from "react";
import { AppContext } from "../Layouts/DefaultLayout";

const ECommerce: React.FC = () => {
  const {user} = useContext(AppContext)
  return (
    <>
      <div className="grid grid-cols-1 gap-4 md:grid-cols-2 md:gap-6 xl:grid-cols-4 2xl:gap-7.5">
        Hello {user?.first_name+" "+user?.last_name},
      </div>
    </>
  );
};

export default ECommerce;
