"use client";
import React, { useState, ReactNode, useEffect, createContext } from "react";
import Sidebar from "@/components/Sidebar";
import Header from "@/components/Header";
import Cookies from 'js-cookie'
import apiInstance, { APIResponse } from "@/services/api";
import { QueryClient, QueryClientProvider, useQuery } from "@tanstack/react-query";
import ms from 'ms'
import User from "@/types/user";
import { AxiosError } from "axios";
import toast from "react-hot-toast";

const queryClient = new QueryClient()
type AppContextType = {
  user: User | undefined
  logOut: () => Promise<void>
}
export const AppContext = createContext<AppContextType>({
  user: undefined,
  logOut: async () => {}
});

function AuthProvider({
  children,
}: {
  
  children: React.ReactNode;
}) {
  const { data: userRes } = useQuery<APIResponse<User>, AxiosError>({
    queryKey: ['user'],
    staleTime: ms("15m"),
    queryFn: () => apiInstance
      .get('/api/v1/user/profile'),
  })
  
  const logOut = async () => {
    toast.loading("Loading...", {id: "loading"})
    Cookies.remove('accessToken');
    try {
      const response = await apiInstance.post(
        '/api/v1/auth/logout',
        {},
        { withCredentials: true }
      );
      toast.success(response.data.message)
      toast.remove("loading")
      window.location.replace("/auth/signin")
    } catch (error) {
      toast.remove("loading")
      if (error instanceof AxiosError) {
        if (!!error.response && error.response?.status >= 400)
          toast.error(error.response?.data?.message)
      } else toast.error("Server error")
    }
  }

  return <AppContext.Provider value={{ user: userRes?.data.data, logOut }}>
     {children}
  </AppContext.Provider>
}

export default function DefaultLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const [sidebarOpen, setSidebarOpen] = useState(false);
  

  return (
      <QueryClientProvider client={queryClient}>
      <AuthProvider>
        {/* <!-- ===== Page Wrapper Start ===== --> */}
        <div className="flex">
          {/* <!-- ===== Sidebar Start ===== --> */}
          <Sidebar sidebarOpen={sidebarOpen} setSidebarOpen={setSidebarOpen} />
          {/* <!-- ===== Sidebar End ===== --> */}

          {/* <!-- ===== Content Area Start ===== --> */}
          <div className="relative flex flex-1 flex-col lg:ml-72.5">
            {/* <!-- ===== Header Start ===== --> */}
            <Header sidebarOpen={sidebarOpen} setSidebarOpen={setSidebarOpen} />
            {/* <!-- ===== Header End ===== --> */}

            {/* <!-- ===== Main Content Start ===== --> */}
            <main>
              <div className="mx-auto max-w-screen-2xl p-4 md:p-6 2xl:p-10">
                {children}
              </div>
            </main>
            {/* <!-- ===== Main Content End ===== --> */}
          </div>
          {/* <!-- ===== Content Area End ===== --> */}
        </div>
        {/* <!-- ===== Page Wrapper End ===== --> */}
      </AuthProvider>
        
      </QueryClientProvider>
   
  );
}
