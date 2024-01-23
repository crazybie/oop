for a in range(8):
    for r in range(5):
        geneArg = ','.join('R%d' % i for i in range(r))
        if a:
            if r:
                geneArg += ','
            geneArg += ','.join('A%d' % i for i in range(a))
        if geneArg:
            geneArg = '[' + geneArg + ' any]'
        fnArgs = ','.join('A%d' % i for i in range(a))
        args = ','.join(f'a{i} A{i}' for i in range(a))
        argNames = ','.join('a%d' % i for i in range(a))
        outArgs = ','.join('checkedCast[R%d](r[%d])' % (i, i) for i in range(r))
        retTypes = ','.join('R%d' % i for i in range(r))
        if r:
            retTypes = '(' + retTypes + ')'

        print(f'''// nolint:lll
func Invoke{a}_{r}{geneArg}(f func({fnArgs}) {retTypes}, obj iStruct {',' if args else ''} {args}) {retTypes}{{
    {'r := ' if r else ''}obj.call(f {',' if argNames else ''} {argNames})
    return {outArgs}
}}''')
