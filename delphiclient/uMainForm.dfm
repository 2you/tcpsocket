object MainForm: TMainForm
  Left = 0
  Top = 0
  Caption = 'MainForm'
  ClientHeight = 410
  ClientWidth = 787
  Color = clBtnFace
  Font.Charset = DEFAULT_CHARSET
  Font.Color = clWindowText
  Font.Height = -11
  Font.Name = 'Tahoma'
  Font.Style = []
  OldCreateOrder = False
  PixelsPerInch = 96
  TextHeight = 13
  object Button1: TButton
    Left = 488
    Top = 30
    Width = 75
    Height = 25
    Caption = #36830#25509'1'
    TabOrder = 0
    OnClick = Button1Click
  end
  object edtIP: TEdit
    Left = 128
    Top = 32
    Width = 121
    Height = 21
    TabOrder = 1
    Text = '127.0.0.1'
  end
  object edtPort: TEdit
    Left = 312
    Top = 32
    Width = 121
    Height = 21
    TabOrder = 2
    Text = '11223'
  end
  object Button2: TButton
    Left = 616
    Top = 30
    Width = 75
    Height = 25
    Caption = #26029#24320'1'
    TabOrder = 3
    OnClick = Button2Click
  end
  object Button3: TButton
    Left = 488
    Top = 72
    Width = 75
    Height = 25
    Caption = #21457#36865'10'
    TabOrder = 4
    OnClick = Button3Click
  end
  object Button4: TButton
    Left = 616
    Top = 72
    Width = 75
    Height = 25
    Caption = #21457#36865'11'
    TabOrder = 5
    OnClick = Button4Click
  end
  object Button5: TButton
    Left = 128
    Top = 134
    Width = 75
    Height = 25
    Caption = #36830#25509'1'
    TabOrder = 6
    OnClick = Button5Click
  end
  object Button6: TButton
    Left = 256
    Top = 134
    Width = 75
    Height = 25
    Caption = #26029#24320'1'
    TabOrder = 7
    OnClick = Button6Click
  end
  object Button7: TButton
    Left = 128
    Top = 176
    Width = 75
    Height = 25
    Caption = #21457#36865'10'
    TabOrder = 8
    OnClick = Button7Click
  end
  object Button8: TButton
    Left = 256
    Top = 176
    Width = 75
    Height = 25
    Caption = #21457#36865'11'
    TabOrder = 9
    OnClick = Button8Click
  end
  object ClientSocket1: TClientSocket
    Active = False
    ClientType = ctNonBlocking
    Port = 0
    Left = 56
    Top = 40
  end
  object ClientSocket2: TClientSocket
    Active = False
    ClientType = ctNonBlocking
    Port = 0
    Left = 368
    Top = 96
  end
end
