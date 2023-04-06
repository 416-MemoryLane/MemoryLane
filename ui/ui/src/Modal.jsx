import { CloseIcon } from "./Icons/CloseIcon";

export const Modal = ({ isOpen, children, callbackFn }) => {
  return isOpen ? (
    <div className="fixed flex h-screen w-screen flex-col content-center items-center justify-center bg-gray-500/50 backdrop-blur-md z-50">
      <div className={`flex max-h-screen  max-w-[1440px] flex-col`}>
        <div className="fixed z-50 -mr-5 place-self-end ">
          <CloseIcon
            className="z-50 rounded-full bg-slate-300 hover:cursor-pointer"
            onClick={callbackFn}
          />
        </div>
        <div className="flex justify-center my-4 rounded-2xl border border-gray-400 bg-white p-10 shadow-2xl">
          {children}
        </div>
      </div>
    </div>
  ) : (
    <></>
  );
};
