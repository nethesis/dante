Name:		dante
Version:	0.0.0
Release:	0%{?dist}
Summary:	Single stack reports made simple

License:	GPLv3
URL:	    https://github.com/nethesis/dante	
Source0:	https://github.com/nethesis/dante/archive/master.tar.gz
Source1:    caronte
Source2:    dante.sysconf


%description
Single stack reports made simple


%prep
%setup -q -n dante-master


%install
mkdir -p %{buildroot}/usr/share/dante/
mkdir -p %{buildroot}/usr/bin
mkdir -p %{buildroot}/etc/sysconfig/
cp ciacco/ciacco %{buildroot}/%{_bindir}
mv %{SOURCE1}  %{buildroot}/%{_bindir}
mv %{SOURCE2}  %{buildroot}/etc/sysconfig/dante


%files
%doc README.md
%license LICENSE
%config /etc/sysconfig/dante
/usr/share/dante/
%{_bindir}/ciacco
%{_bindir}/caronte



%changelog

