for a in range(5):
    for r in range(3):
        geneArg = ','.join('R%d'%i for i in range(r))  
        if a:
            if r:
                geneArg += ','
            geneArg += ','.join('A%d'%i for i in range(a))
        if geneArg:
            geneArg = '[' + geneArg + ' any]'
        fnArgs = ','.join('A%d'%i for i in range(a))
        args = ','.join(f'a{i} A{i}' for i in range(a))
        argNames = ','.join('a%d'%i for i in range(a))
        outArgs = ','.join('out[%d].Interface().(R%d)'%(i,i) for i in range(r))
        retTypes = ','.join('R%d'%i for i in range(r))
        if r:
            retTypes = '(' + retTypes + ')'
        print(f'''func CallArg{a}Ret{r}{geneArg}(obj iStruct, f func({fnArgs}) {retTypes} {',' if args else ''} {args}) {retTypes}{{
    {'out := ' if r else ''}obj.Call(f {',' if argNames else ''} {argNames})
    return {outArgs}
}}''')