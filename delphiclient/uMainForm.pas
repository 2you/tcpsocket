unit uMainForm;

interface

uses
  Winapi.Windows, Winapi.Messages, System.SysUtils, System.Variants, System.Classes, Vcl.Graphics,
  Vcl.Controls, Vcl.Forms, Vcl.Dialogs, System.Win.ScktComp, Vcl.StdCtrls, IdGlobal;

type
  TMainForm = class(TForm)
    ClientSocket1: TClientSocket;
    Button1: TButton;
    edtIP: TEdit;
    edtPort: TEdit;
    Button2: TButton;
    Button3: TButton;
    Button4: TButton;
    Button5: TButton;
    Button6: TButton;
    Button7: TButton;
    Button8: TButton;
    ClientSocket2: TClientSocket;
    procedure Button1Click(Sender: TObject);
    procedure Button2Click(Sender: TObject);
    procedure Button3Click(Sender: TObject);
    procedure Button4Click(Sender: TObject);
    procedure Button5Click(Sender: TObject);
    procedure Button6Click(Sender: TObject);
    procedure Button7Click(Sender: TObject);
    procedure Button8Click(Sender: TObject);
  private
    { Private declarations }
  public
    { Public declarations }
  end;

var
  MainForm: TMainForm;

implementation

{$R *.dfm}

procedure TMainForm.Button1Click(Sender: TObject);
begin
  ClientSocket1.Close;
  ClientSocket1.Host := edtIP.Text;
  ClientSocket1.Port := StrToInt(edtPort.Text);
  ClientSocket1.Open;
end;

procedure TMainForm.Button2Click(Sender: TObject);
begin
  ClientSocket1.Close;
end;

procedure TMainForm.Button3Click(Sender: TObject);
var
  sText: string;
  nLen: Integer;
  buf1: TIdBytes;
begin
  sText := 'abcdefºº×Ö';
  buf1 := ToBytes(sText, IndyTextEncoding_UTF8);
  nLen := Length(buf1);
  ClientSocket1.Socket.SendBuf(nLen, 4);
  ClientSocket1.Socket.SendBuf(buf1[0], nLen);
end;

procedure TMainForm.Button4Click(Sender: TObject);
var
  sText: string;
  nLen: Integer;
  buf1: TIdBytes;
begin
  sText := '12345';
  buf1 := ToBytes(sText, IndyTextEncoding_UTF8);
  nLen := Length(buf1);
  ClientSocket1.Socket.SendBuf(nLen, 4);
  ClientSocket1.Socket.SendBuf(buf1[0], nLen);
end;

procedure TMainForm.Button5Click(Sender: TObject);
begin
  ClientSocket2.Close;
  ClientSocket2.Host := edtIP.Text;
  ClientSocket2.Port := StrToInt(edtPort.Text);
  ClientSocket2.Open;
end;

procedure TMainForm.Button6Click(Sender: TObject);
begin
  ClientSocket2.Close;
end;

procedure TMainForm.Button7Click(Sender: TObject);
var
  sText: string;
  nLen: Integer;
  buf1: TIdBytes;
begin
  sText := 'abcdefºº×Ö';
  buf1 := ToBytes(sText, IndyTextEncoding_UTF8);
  nLen := Length(buf1);
  ClientSocket2.Socket.SendBuf(nLen, 4);
  ClientSocket2.Socket.SendBuf(buf1[0], nLen);
end;

procedure TMainForm.Button8Click(Sender: TObject);
var
  sText: string;
  nLen: Integer;
  buf1: TIdBytes;
begin
  sText := '12345';
  buf1 := ToBytes(sText, IndyTextEncoding_UTF8);
  nLen := Length(buf1);
  ClientSocket2.Socket.SendBuf(nLen, 4);
  ClientSocket2.Socket.SendBuf(buf1[0], nLen);
end;

end.
