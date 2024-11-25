"use client";
import { AppContext } from '@/components/Layouts/DefaultLayout'
import apiInstance from '@/services/api'
import { AxiosError } from 'axios'
import Image from 'next/image'
import React, { useCallback, useContext, useState } from 'react'
import toast from 'react-hot-toast'
import { useDropzone } from 'react-dropzone'
import { AlertCircle, CheckCircle2, FileText, ImageIcon, X } from 'lucide-react'
import { Progress } from '@/components/ui/progress'

const EditPhotoForm = () => {
    const [image, setImage] = useState<string>("")
    const [file, selectFile] = useState<File>()
    const [uploadProgress, setUploadProgress] = useState(0)
    const [uploadError, setUploadError] = useState<string | null>(null)
    const { user } = useContext(AppContext)

    const onDrop = useCallback((acceptedFiles: File[]) => {
        const file = acceptedFiles[0]
        if (file) {
            selectFile(file)

            setUploadError(null)
            const reader = new FileReader();
            reader.onload = () => {
                setImage(!!reader.result? reader.result as string: "");
            };
            reader.readAsDataURL(file);
            // Simulate upload progress
            let progress = 0
            const interval = setInterval(() => {
                progress += 10
                setUploadProgress(progress)
                if (progress >= 100) {
                    clearInterval(interval)
                }
            }, 200)
        } else {
            setUploadError('Please upload a valid PDF file.')
        }
    }, [])

    const { getRootProps, getInputProps, isDragActive } = useDropzone({
        onDrop,
        accept: { 'image/jpeg': ['.jpg', ".jpeg"], 'image/png': ['.png'], 'image/gif': [".gif"] },
        multiple: false
    })

    const onSubmit = async () => {
        try {
            if(!file){
                toast.error("Please select file")
                return
            }
            const formData = new FormData()
            formData.append("image", file)
            toast.loading("Loading...", { id: "loading" })
            const response = await apiInstance.put('/api/v1/user/photo', formData);

            toast.success(response.data.message)
            toast.remove("loading")
        } catch (error) {
            setUploadError('Upload error')
            toast.remove("loading")
            if (error instanceof AxiosError) {
                if (!!error.response && error.response?.status >= 400)
                    toast.error(error.response?.data?.message)
            } else toast.error("Server error")
        }
    }
  return (
      <form action="#">
          <div className="mb-4 flex items-center gap-3">
              <div className="h-14 w-14 rounded-full">
                  {!!user?.photo && !image.length && <Image
                      src={user.photo}
                      width={55}
                      height={55}
                      alt="User"
                  />}

                  {!user?.photo && !image.length && <ImageIcon color='black' size={55} />}
                  {!!image.length && <Image
                      src={image}
                      width={55}
                      height={55}
                      className='rounded-3xl'
                      alt="User"
                  />}
              </div>
              <div>
                  <span className="mb-1.5 text-black dark:text-white">
                      Edit your photo
                  </span>
                  
              </div>
          </div>

          {!file ? (
              <div
                  {...getRootProps()}
                  className={`border-2 mb-6 border-dashed rounded-lg p-8 text-center cursor-pointer transition-colors ${isDragActive ? 'border-primary bg-primary/10' : 'border-gray-300 hover:border-primary'
                      }`}
              >
                  <input {...getInputProps()} />
                  <FileText className="mx-auto h-12 w-12 text-gray-400" />
                  <p className="mt-2 text-sm text-gray-600">
                      Drag & drop your image here, or click to select a file
                  </p>
              </div>
          ) : (
                  <div className="space-y-4 mb-6">
                  <div className="flex items-center space-x-2">
                      <FileText className="h-6 w-6 text-blue-500" />
                      <span className="font-medium">{file.name}</span>
                      <button
                          
                         
                          className="ml-auto"
                          onClick={() => {
                            selectFile(undefined)
                          }}
                      >
                          <X className="h-4 w-4" />
                          <span className="sr-only">Remove file</span>
                      </button>
                  </div>
                  <Progress value={uploadProgress} className="w-full" />
                  {uploadProgress === 100 && (
                      <div className="flex items-center text-green-600">
                          <CheckCircle2 className="mr-2 h-4 w-4" />
                          <span>Upload complete!</span>
                      </div>
                  )}
              </div>
          )}
          {uploadError && (
              <div className="flex items-center text-red-600 mt-2">
                  <AlertCircle className="mr-2 h-4 w-4" />
                  <span>{uploadError}</span>
              </div>
          )}

          <div className="flex justify-end gap-4.5">
              <button
                  className="flex justify-center rounded border border-stroke px-6 py-2 font-medium text-black hover:shadow-1 dark:border-strokedark dark:text-white"
                  onClick={(e) => {
                    e.preventDefault()
                    selectFile(undefined)
                }}
              >
                  Cancel
              </button>
              <button
                  onClick={e => {
                    e.preventDefault()
                    onSubmit()}}
                  className="flex justify-center rounded bg-primary px-6 py-2 font-medium text-gray hover:bg-opacity-90"
                  type="submit"
              >
                  Save
              </button>
          </div>
      </form>
  )
}

export default EditPhotoForm