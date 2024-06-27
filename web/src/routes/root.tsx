import { MutableRefObject, useEffect, useRef, useState } from "react";

const baseApiUrl = "localhost:8000/api/v1";
type InvoiceId = string;

interface InvoiceResponse {
  id: InvoiceId;
}

enum InvoiceStatus {
  Uploaded,
  Processed,
}

interface Invoice {
  id: InvoiceId;
  name: string;
  status: InvoiceStatus;
}

export default function Root() {
  const [inputValue, setInputValue] = useState("");
  const [fileTableBody, setFileTableBody] = useState<React.JSX.Element[]>([]);
  const [files, setFiles] = useState<Invoice[]>([]);

  const updateFiles = (files: Invoice[]) => {
    setFiles(files);
    setFileTableBody(
      files.map((file) => {
        return (
          <tr
            className="bg-white border-b dark:bg-gray-800 dark:border-gray-700"
            key={file.id}
          >
            <th
              scope="row"
              className="px-6 py-4 font-medium text-gray-900 whitespace-nowrap dark:text-white"
            >
              {file.id}
            </th>
            <td className="px-6 py-4">{file.name}</td>
            <td className="px-6 py-4">{InvoiceStatus[file.status]}</td>
          </tr>
        );
      }),
    );
  };

  const handleInputValueChange = (
    event: React.ChangeEvent<HTMLInputElement>,
  ) => {
    setInputValue(event.target.value);
  };

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    const request = new Request(`//${baseApiUrl}/invoices`, {
      body: new FormData(event.currentTarget),
      method: "POST",
    });

    const response = await fetch(request);
    const body: InvoiceResponse = await response.json();
    files.push({
      id: body.id,
      // TODO: This is probably only valid for UNIX systems
      name: inputValue.split("\\").at(-1) as string,
      status: InvoiceStatus.Uploaded,
    });
    updateFiles(files);
  };

  const webSocket: MutableRefObject<WebSocket | null> = useRef(null);
  useEffect(() => {
    webSocket.current = new WebSocket(`ws://${baseApiUrl}/websocket`);

    webSocket.current.onerror = (event) => {
      console.log(event);
    };

    webSocket.current.onmessage = (event) => {
      const parsedData: {
        type: string;
        value: InvoiceId;
      } = JSON.parse(event.data);
      files.map((file) => {
        if (file.id === parsedData.value) {
          file.status = InvoiceStatus.Processed;
        }

        return file;
      });
      updateFiles(files);
    };

    return () => {
      const currentWs = webSocket.current;
      if (currentWs === null) {
        return;
      }

      currentWs.close();
    };
  }, [files]);

  return (
    <main className="flex min-h-screen flex-col items-center p-24">
      <form onSubmit={handleSubmit}>
        <label htmlFor="file" className="block text-sm font-medium leading-6">
          File *
        </label>
        <input
          type="file"
          name="file"
          value={inputValue}
          onChange={handleInputValueChange}
        />

        <button type="submit" className="block text-sm font-medium leading-6">
          Submit
        </button>
      </form>

      <div className="relative overflow-x-auto w-full mt-8">
        <table className="w-full text-sm text-left rtl:text-right text-gray-500 dark:text-gray-400">
          <thead className="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400">
            <tr>
              <th scope="col" className="px-6 py-3">
                ID
              </th>
              <th scope="col" className="px-6 py-3">
                Filename
              </th>
              <th scope="col" className="px-6 py-3">
                Status
              </th>
            </tr>
          </thead>
          <tbody>{fileTableBody}</tbody>
        </table>
      </div>
    </main>
  );
}
