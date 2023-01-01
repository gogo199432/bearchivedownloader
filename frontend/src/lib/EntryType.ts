export class Entry{
    Id :string
    Url  :string
    Title        :string
    Text         : string
    Date         : Date
    Author       :string
    Tags         :string[]

    Choices     :Record<string, string>

    constructor(Id: string, Url: string, Title: string, Text: string, Date: Date, Author: string, Tags: string[], Choices: Record<string, string>) {
        this.Url = Url;
        this.Title = Title;
        this.Text = Text;
        this.Date = Date;
        this.Author = Author;
        this.Tags = Tags;
        this.Choices = Choices
        this.Id = Id
    }
}