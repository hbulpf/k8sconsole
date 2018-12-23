import { Component, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';

@Component({
  selector: 'app-namespace',
  templateUrl: './namespace.component.html',
  styleUrls: ['./namespace.component.css']
})
export class NamespaceComponent implements OnInit {

  // tslint:disable-next-line:no-inferrable-types
  public anyList: any;
  public f: any = 'hello';

  constructor(private http: HttpClient) { }

  ngOnInit() {
    this.http.get('http://50125.hnbdata.cn:63453/api/v1/node?filterBy=name%2Cminikube&sortBy=d%2Cname&itemsPerPage=1&page=1')
    .subscribe(res => { this.anyList = res; });
  }

}
