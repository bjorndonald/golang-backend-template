
import React, {  } from "react";

export default function AuthLayout({
    children,
}: {
    children: React.ReactNode;
}) {
    return (
        <>
            <div className="flex">
                <div className="relative flex flex-1 flex-col lg:max-w-6xl mx-auto px-4">
                    
                    <main>
                        <div className="mx-auto max-w-screen-2xl p-4 md:p-6 2xl:p-10">
                            {children}
                        </div>
                    </main>
                    
                </div>
               
            </div>
        </>
    );
}
